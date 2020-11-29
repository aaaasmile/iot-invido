export default {
    state: {
        timeStamp: '',
        tempraw: '',
        press: '',
        humiraw: '',
        gaso: '',
        iaq: '',
        iaqacc: '',
        temp: '',
        humy: '',
        co2: '',
        voc: '',
        sensorid: '',
        place: '',
    },
    mutations: {
        sensorstate(state, dataArr) {
            if (dataArr.length <= 0){
                return
            }
            let data = dataArr[dataArr.length - 1]
            state.timeStamp = data.timeStamp
            state.tempraw = data.tempraw
            state.press = data.press
            state.humiraw = data.humiraw
            state.gaso = data.gaso
            state.iaq = data.iaq
            state.iaqacc = data.iaqacc
            state.temp = data.temp
            state.humy = data.humy
            state.co2 = data.co2
            state.voc = data.voc
            state.sensorid = data.sensorid
            state.place = data.place
        }
    }
}