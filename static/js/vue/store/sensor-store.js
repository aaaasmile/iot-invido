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
    },
    mutations: {
        sensorstate(state, data) {
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

        }
    }
}