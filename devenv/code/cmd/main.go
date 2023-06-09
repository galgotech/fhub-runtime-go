// Copyright 2023 The fhub-runtime-go Authors
// This file is part of fhub-runtime-go.
//
// This file is part of fhub-runtime-go.
// fhub-runtime-go is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// fhub-runtime-go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with fhub-runtime-go. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"log"

	fhub "github.com/galgotech/fhub-go"
	"github.com/galgotech/fhub-go/devenv/code/pkg"
)

func main() {
	f := &pkg.Functions{}
	fhub.SetPath("devenv/code")
	err := fhub.Run(f)
	if err != nil {
		log.Fatal(err)
	}
}
