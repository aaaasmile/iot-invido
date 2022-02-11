
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
			that.$store.commit('sensorstate', result.data.dataview)
			that.loading = false
		}, error => {
			handleError(error, that)
		});
	},
	InsertData(that, req) {
		console.log('Request is ', req)
		that.$http.post("InsertTestData", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.$store.commit('sensorstate', result.data.dataview)
			that.loading = false
		}, error => {
			handleError(error, that)
		});
	},
	CheckAPIToken(that, req, fnOK){
		console.log('Request is ', req)
		that.$http.post("CheckAPIToken", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.loading = false
      if (fnOK){
        fnOK(result.data)
      }
		}, error => {
			handleError(error, that)
		});
	},
  SignIn(that, req, fnOK){
		console.log('Request is ', req)
		that.$http.post("SignIn", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
			console.log('Call result ', result.data)
			that.loading = false
			that.$store.commit('clearMsgText')
      if (fnOK){
        fnOK(result.data)
      }
		}, error => {
			handleError(error, that)
		});
	},
}