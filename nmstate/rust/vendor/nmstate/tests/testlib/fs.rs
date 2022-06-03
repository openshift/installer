pub(crate) fn read_folder(folder_path: &str) -> Vec<String> {
    let mut folder_contents = Vec::new();
    for entry in std::fs::read_dir(folder_path).unwrap() {
        let entry = entry.unwrap();
        let path = entry.path();
        folder_contents.push(
            path.strip_prefix(folder_path)
                .unwrap()
                .to_str()
                .unwrap()
                .to_string(),
        );
    }
    folder_contents
}
