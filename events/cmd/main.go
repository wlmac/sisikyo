// This file is part of a program named Sisikyō or Sisikyo.
//
// Copyright (C) 2019 Ken Shibata <kenxshibata@gmail.com>
//
// License as published by the Free Software Foundation, either version 1 of the License, or (at your option) any later
// version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
// warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see
// <https://www.gnu.org/licenses/>.

// Package cmd provides common facilities for commands.
package cmd

import (
	_ "embed"
)

const StartupInfo = `This program, named Sisikyō or Sisikyo, is a program that fetches event information from an API.
Copyright (C) 2021 Ken Shibata <kenxshibata@gmail.com>
This program comes with ABSOLUTELY NO WARRANTY and this program is free software, and you are welcome to redistribute it under certain conditions; for details run with the '-h' flag or view the 'license.md' file. This program uses open source software; for details run with the '-h' flag or view the 'license.md' file.
`

//go:embed license_info.md
var LicenseInfo string
