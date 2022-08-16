// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <alex@unikraft.io>
//
// Copyright (c) 2022, Unikraft GmbH. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package processtree

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Render
	red        = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render
	blue       = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Render
	green      = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render
	logStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("254")).Render
	lightGrey  = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render
)

func (pt ProcessTree) View() string {
	s := ""

	finished := 0

	// Update timers on all active items and their parents
	TraverseTreeAndCall(pt.tree, func(pti *ProcessTreeItem) error {
		if pti.status == StatusSuccess ||
			pti.status == StatusFailed ||
			pti.status == StatusFailedChild {
			finished += 1
		}

		return nil
	})

	for _, pti := range pt.tree {
		s += pt.printItem(pti, 0)
	}

	if len(pt.verb) > 0 {
		title := titleStyle(
			pt.verb + " " + formatTimer(pt.timer.Elapsed()) +
				" (" + strconv.Itoa(finished) + "/" + strconv.Itoa(pt.total) + ")",
		)
		s = title + "\n" + s
	}

	if !pt.quitting {
		s += lightGrey("ctrl+c to cancel\n")
	}

	return s
}

func (stm ProcessTree) printItem(pti *ProcessTreeItem, offset uint) string {
	failed := 0
	completed := 0
	running := 0

	// Determine the status of immediate children whilst we're printing
	for _, child := range pti.children {
		if child.status == StatusFailed ||
			child.status == StatusFailedChild {
			failed += 1
		} else if child.status == StatusSuccess {
			completed += 1
		} else if child.status == StatusRunningChild ||
			child.status == StatusRunning ||
			child.status == StatusRunningButAChildHasFailed {
			running += 1
		}
	}

	if len(pti.children) > 0 {
		if failed > 0 && running > 0 {
			pti.status = StatusRunningButAChildHasFailed
		} else if running > 0 {
			pti.status = StatusRunningChild
		} else if failed > 0 {
			pti.status = StatusFailedChild
		}
	}

	width := lipgloss.Width

	textLeft := ""
	switch pti.status {
	case StatusSuccess:
		textLeft += green("[+]")
	case StatusFailed, StatusFailedChild:
		textLeft += red("<!>")
	case StatusRunning, StatusRunningChild, StatusRunningButAChildHasFailed:
		textLeft += "(" + pti.spinner.View() + ")"
	default:
		textLeft += " "
	}

	textLeft += " " + pti.textLeft

	elapsed := formatTimer(pti.timer.Elapsed())
	rightTimerWidth := width(elapsed)
	if rightTimerWidth > stm.rightPad {
		stm.rightPad = rightTimerWidth
	}

	textRight := ""
	if len(pti.textRight) > 0 {
		switch pti.status {
		case StatusSuccess:
			textRight += green(pti.textRight)
		case StatusFailed, StatusFailedChild:
			textRight += red(pti.textRight)
		case StatusRunning, StatusRunningChild, StatusRunningButAChildHasFailed:
			textRight += blue(pti.textRight)
		default:
			textRight += pti.textRight
		}
	}

	textRight += lightGrey(" [" +
		lipgloss.NewStyle().
			Render(indent.String(elapsed, uint(stm.rightPad-rightTimerWidth))) +
		"]")

	left := lipgloss.NewStyle().
		Width(stm.width - width(textRight) - int(offset*INDENTS)).
		Height(1).
		Render(textLeft)

	right := lipgloss.NewStyle().
		Width(width(textRight)).
		Height(1).
		Render(textRight)

	s := lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		right,
	) + "\n"

	// Print the logs for this item
	truncate := 0
	loglen := len(pti.logs) - LOGLEN
	if loglen > 0 {
		truncate = loglen
	}
	if pti.status != StatusSuccess {
		for _, line := range pti.logs[truncate:] {
			s += indent.String(logStyle(line), INDENTS) + "\n"
		}
	}

	// Print the child processes
	for _, child := range pti.children {
		s += stm.printItem(child, offset+1)
	}

	// Do not indent the root node
	if offset == 0 {
		return s
	}

	// Since this method is recursive, indent by 1 factor
	return indent.String(s, INDENTS)
}
