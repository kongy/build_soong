// Copyright 2023 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rust

import (
	"android/soong/android"
	"testing"

	"github.com/google/blueprint"
)

var GlobalThinLTOPreparer = android.GroupFixturePreparers(
	prepareForRustTest,
	android.FixtureModifyEnv(func(env map[string]string) {
		env["GLOBAL_THINLTO"] = "false"
	}))

func TestStaticDeps(t *testing.T) {
	bp := `
		rust_library {
			name: "libfoo",
			srcs: ["foo.rs"],
			crate_name: "foo",
		}
		cc_library_static {
			name: "libbar",
			srcs: ["bar.cpp"]
		}
		`

	result := GlobalThinLTOPreparer.RunTestWithBp(t, bp)

	hasDep := func(m android.Module, wantDep android.Module) bool {
		var found bool
		result.VisitDirectDeps(m, func(dep blueprint.Module) {
			if dep == wantDep {
				found = true
			}
		})
		return found
	}

	libfoo := result.ModuleForTests("libfoo", "android_arm64_armv8-a_dylib").Module()
	libbar := result.ModuleForTests("libbar", "android_arm64_armv8-a_static_lto-none").Module()
	if !hasDep(libfoo, libbar) {
		t.Errorf("'libfoo' missing dependency on no lto variant of 'libbar'")
	}
}
