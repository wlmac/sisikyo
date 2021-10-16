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

// Package api provides abstractions for the REST API.
package api

import (
	"net/http"
)

// HTTPClient is a subset of the *http.Client interface.
type HTTPClient interface {
	Do(r *http.Request) (resp *http.Response, err error)
}

type Req interface {
	Req(c *Client) (*http.Request, error)
}

var Reqs = [...]Req{
	AuthReq{},
	AuthReReq{},
	AnnFeedReq{},
	AnnReq{},
	OrgReq{},
	OrgsReq{},
	UserReq{},
	MeReq{},
	MeScheduleReq{},
	MeScheduleWeekReq{},
	MeTimetableReq{},
	EventsReq{},
	//TimetablesReq{},
	//TimetableScheduleReq{},
	//TimetableReq{},
	TermsReq{},
	//TermReq{},
	//TermCurrentReq{},
	//TermScheduleReq{},
	//TermScheduleWeekReq{},
	// martor_image_upload unsupported
	VersionReq{},
}
