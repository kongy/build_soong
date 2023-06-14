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
	"android/soong/cc"
)

// Propagate lto requirements down from binaries
func ltoDepsMutator(mctx android.TopDownMutatorContext) {
	if _, ok := mctx.Module().(*Module); ok {
		mctx.WalkDeps(func(dep android.Module, parent android.Module) bool {
			tag := mctx.OtherModuleDependencyTag(dep)
			if !cc.IsStaticDepTag(tag) {
				return false
			}

			if dep, ok := dep.(*cc.Module); ok {
				dep.SetLtoIsRustDep(true)
			}

			return true
		})
	}
}

func ltoMutator(mctx android.BottomUpMutatorContext) {
	if cc.GlobalThinLTO(mctx) {
		mctx.SetDependencyVariation("lto-none")
	}
}
