use std::io;
use std::thread;
use std::time::Duration;

use easy_parallel::Parallel;
use polling::Poller;

#[test]
fn simple() -> io::Result<()> {
    let poller = Poller::new()?;
    let mut events = Vec::new();

    for _ in 0..10 {
        poller.notify()?;
        poller.wait(&mut events, None)?;
    }

    Ok(())
}

#[test]
fn concurrent() -> io::Result<()> {
    let poller = Poller::new()?;
    let mut events = Vec::new();

    for _ in 0..2 {
        Parallel::new()
            .add(|| {
                thread::sleep(Duration::from_secs(0));
                poller.notify().unwrap();
            })
            .finish(|| poller.wait(&mut events, None).unwrap());
    }

    Ok(())
}
