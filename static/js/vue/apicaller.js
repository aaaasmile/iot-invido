
const handleError = (error, that) => {
	console.error(error);
	that.loading = false
	if (error.bodyText !== '') {
		that.$store.commit('msgText', `${error.statusText}: ${error.bodyText}`)
	} else {
		that.$store.commit('msgText', 'Error: empty response')
	}
}

export default {
	FetchData(that, req) {
		console.log('Request is ', req)
		that.$http.post("FetchData", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('sensorstate', result.data)
			that.loading = false
		}, error => {
			handleError(error, that)
		});
	}
}