// Copyright © 2018 The Havener
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package havener

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var FixedTerminalWidth = -1
var FixedTerminalHeight = -1

var DefaultTerminalWidth = 80
var DefaultTerminalHeight = 25

func GetTerminalWidth() int {
	width, _ := GetTerminalSize()
	return width
}

func GetTerminalHeight() int {
	_, height := GetTerminalSize()
	return height
}

func GetTerminalSize() (int, int) {
	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))

	switch {
	case err != nil: // Return default fall-back value
		return DefaultTerminalWidth, DefaultTerminalHeight

	case FixedTerminalWidth > 0 && FixedTerminalHeight > 0: // Return user preference (explicit overwrite)
		return FixedTerminalWidth, FixedTerminalHeight

	case FixedTerminalWidth > 0:
		return FixedTerminalWidth, height

	case FixedTerminalHeight > 0:
		return width, FixedTerminalHeight

	default:
		return width, height
	}
}
