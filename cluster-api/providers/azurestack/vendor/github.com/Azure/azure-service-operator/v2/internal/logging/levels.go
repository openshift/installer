/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package logging

// These levels are defined to categorize (roughly) what a given log level means. They serve as guidelines
// to help decide what level a new log should be emitted at. They are not rules that must be followed.
// Keep in mind that there may be logs logged at a level higher
// than Debug (the highest log level defined below). Logs between levels 4 and 10 are progressively
// more verbose and are used rarely.

const (
	// Status logs are logs that should ALWAYS be shown. In the context of a controller,
	// Status logs should not be happening for every reconcile loop. Examples of good Status logs
	// are state changes or major problems. These probably line up very nicely with interesting
	// customer facing "events" as well.
	Status = 0

	// Info logs are logs that are probably useful but may be slightly more verbose.
	// In the context of a controller, an Info log probably shouldn't be emitted every time through
	// the reconcile loop, at least in the happy path.
	// Examples of good Info logs include intermittent errors which we expect to be able to retry through
	// and object updates (think: "set ownerReference" or similar things which are not critical state changes
	// but are still interesting updates that aren't super verbose).
	Info = 1

	// Verbose logs are logs that are quite verbose. In the context of a controller
	// they likely log on each reconcile loop. Examples of good Verbose logs include
	// "waiting for deployment..." or "waiting for deletion to complete..."
	Verbose = 2

	// Debug logs are logs that are extremely verbose and log each reconcile loop (or multiple times in a single
	// reconcile loop). Examples include ARM request and response payloads, or request and response payloads
	// from APIServer.
	Debug = 3
)
