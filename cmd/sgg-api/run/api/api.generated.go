package api

//
// ATTENTION: This file is generated automagically.
// Do not touch it. Do not pass go. Do not collect $200.
// Instead run 'go generate' or 'make gen' to build this file.
//

import (
	mw "save.gg/sgg/cmd/sgg-api/run/middleware"
	"save.gg/sgg/meta"
)

func init() {

	meta.RegisterRoute("GET", "/api/user/:slug",
		mw.VR(mw.VRMap{
			"default": mw.SC(getUser),

			"v1": mw.SC(getUser),
		}),
	)

	meta.RegisterRoute("PATCH", "/api/user/:slug",
		mw.VR(mw.VRMap{
			"default": mw.RequireSession(patchUser,
				&mw.SecurityFlags{
					All: true,
				}),

			"v1": mw.RequireSession(patchUser,
				&mw.SecurityFlags{
					All: true,
				}),
		}),
	)

}
