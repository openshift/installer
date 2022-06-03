// Copyright 2021 Red Hat, Inc.
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

use std::process::Command;

pub fn clear_network_environment() {
    cmd_exec("../../tools/test_env", vec!["rm"]);
}

pub fn set_network_environment(env_type: &str) {
    assert!(cmd_exec("../../tools/test_env", vec![env_type]));
}

pub fn cmd_exec(command: &str, args: Vec<&str>) -> bool {
    let mut proc = Command::new(command);
    for argument in args.iter() {
        proc.arg(argument);
    }
    let status = proc.status().expect("failed to execute the command");

    return status.success();
}
