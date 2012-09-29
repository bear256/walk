// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type WebView struct {
	Widget        **walk.WebView
	Name          string
	StretchFactor int
	Row           int
	RowSpan       int
	Column        int
	ColumnSpan    int
	URL           string
}

func (wv WebView) Create(parent walk.Container) error {
	w, err := walk.NewWebView(parent)
	if err != nil {
		return err
	}

	return InitWidget(wv, w, func() error {
		if err := w.SetURL(wv.URL); err != nil {
			return err
		}

		if wv.Widget != nil {
			*wv.Widget = w
		}

		return nil
	})
}

func (wv WebView) CommonInfo() (name string, stretchFactor, row, rowSpan, column, columnSpan int) {
	return wv.Name, wv.StretchFactor, wv.Row, wv.RowSpan, wv.Column, wv.ColumnSpan
}