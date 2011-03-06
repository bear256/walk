// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
	"syscall"
	"unsafe"
)

import (
	. "walk/winapi"
	. "walk/winapi/user32"
)

var comboBoxOrigWndProcPtr uintptr
var _ subclassedWidget = &ComboBox{}

type ComboBox struct {
	WidgetBase
	items                         *ComboBoxItemList
	prevSelIndex                  int
	selectedIndexChangedPublisher EventPublisher
}

func NewComboBox(parent Container) (*ComboBox, os.Error) {
	cb := &ComboBox{prevSelIndex: -1}

	if err := initChildWidget(
		cb,
		parent,
		"COMBOBOX",
		WS_TABSTOP|WS_VISIBLE|WS_VSCROLL|CBS_DROPDOWNLIST,
		0); err != nil {
		return nil, err
	}

	cb.items = newComboBoxItemList(cb)

	return cb, nil
}

func (*ComboBox) origWndProcPtr() uintptr {
	return comboBoxOrigWndProcPtr
}

func (*ComboBox) setOrigWndProcPtr(ptr uintptr) {
	comboBoxOrigWndProcPtr = ptr
}

func (*ComboBox) LayoutFlags() LayoutFlags {
	return ShrinkableHorz | GrowableHorz
}

func (cb *ComboBox) PreferredSize() Size {
	return cb.dialogBaseUnitsToPixels(Size{50, 12})
}

func (cb *ComboBox) Items() *ComboBoxItemList {
	return cb.items
}

func (cb *ComboBox) SelectedIndex() int {
	return int(SendMessage(cb.hWnd, CB_GETCURSEL, 0, 0))
}

func (cb *ComboBox) SetSelectedIndex(value int) os.Error {
	index := int(SendMessage(cb.hWnd, CB_SETCURSEL, uintptr(value), 0))

	if index != value {
		return newError("invalid index")
	}

	if value != cb.prevSelIndex {
		cb.prevSelIndex = value
		cb.selectedIndexChangedPublisher.Publish()
	}

	return nil
}

func (cb *ComboBox) SelectedIndexChanged() *Event {
	return cb.selectedIndexChangedPublisher.Event()
}

func (cb *ComboBox) TextSelection() (start, end int) {
	SendMessage(cb.hWnd, CB_GETEDITSEL, uintptr(unsafe.Pointer(&start)), uintptr(unsafe.Pointer(&end)))
	return
}

func (cb *ComboBox) SetTextSelection(start, end int) {
	SendMessage(cb.hWnd, CB_SETEDITSEL, 0, uintptr(MAKELONG(uint16(start), uint16(end))))
}

func (cb *ComboBox) wndProc(hwnd HWND, msg uint, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		switch HIWORD(uint(wParam)) {
		case CBN_SELENDOK:
			if selIndex := cb.SelectedIndex(); selIndex != cb.prevSelIndex {
				cb.selectedIndexChangedPublisher.Publish()
				cb.prevSelIndex = selIndex
				return 0
			}
		}
	}

	return cb.WidgetBase.wndProc(hwnd, msg, wParam, lParam)
}

func (cb *ComboBox) onInsertingComboBoxItem(index int, item *ComboBoxItem) (err os.Error) {
	if CB_ERR == SendMessage(cb.hWnd, CB_INSERTSTRING, uintptr(index), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(item.text)))) {
		err = newError("CB_INSERTSTRING failed")
	}

	return
}

func (cb *ComboBox) onRemovingComboBoxItem(index int, item *ComboBoxItem) (err os.Error) {
	if CB_ERR == SendMessage(cb.hWnd, CB_DELETESTRING, uintptr(index), 0) {
		err = newError("CB_DELETESTRING failed")
	}

	if index == cb.prevSelIndex {
		cb.prevSelIndex = -1
	}

	return
}

func (cb *ComboBox) onClearingComboBoxItems() (err os.Error) {
	SendMessage(cb.hWnd, CB_RESETCONTENT, 0, 0)

	cb.prevSelIndex = -1

	return
}