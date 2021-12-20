addEventListener('fetch', (event) => {
	event.respondWith(handleRequest(event.request))
})

/**
 * Sample test function
 * Log a given request object
 * @param {Request} request
 */
async function handleRequest(request) {
	console.log('Got request', request)
	const response = await fetch(request)
	return response;
}