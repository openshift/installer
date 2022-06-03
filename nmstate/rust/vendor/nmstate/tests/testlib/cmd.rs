pub(crate) fn cmd_exec_check(cmds: &[&str]) -> String {
    let cmds = cmds.to_vec();
    println!("Executing command {:?}", cmds);
    let run = std::process::Command::new(cmds[0])
        .args(&cmds[1..])
        .output()
        .expect("failed to execute file command");

    let stdout = String::from_utf8(run.stdout)
        .expect("Failed to convert file command output to String");
    println!("STDOUT: {}", &stdout);
    let stderr = String::from_utf8(run.stderr)
        .expect("Failed to convert file command stderr to String");
    if !stderr.is_empty() {
        println!("STDERR: {}", &stderr);
    }
    stdout
}
