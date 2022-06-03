pub(crate) fn with_clean_up_afterwords<T, C>(test: T, cleanup: C)
where
    T: FnOnce() + std::panic::UnwindSafe,
    C: FnOnce(),
{
    let result = std::panic::catch_unwind(|| {
        test();
    });

    cleanup();
    assert!(result.is_ok())
}
