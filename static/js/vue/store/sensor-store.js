export default {
    state: {
        last: {
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
            color: '',
        },
        getClassColor: (iaqclass) => {
            switch (iaqclass) {
                case "Hazardous":
                    return "red"
                case "Very Unhealthy":
                    return "red lighten-4"
                case "More than Unhealthy":
                    return "orange"
                case "Unhealthy":
                    return "amber"
                case "Unhealthy for Sensitive Groups":
                    return "yellow"
                case "Moderate":
                    return "lime accent-3"
                case "Good":
                    return "green"
            }
            return ""
        },
        measures: [],
    },
    mutations: {
        sensorstate(state, dataArr) {
            state.measures = dataArr
            if (dataArr.length <= 0) {
                state.last = {}
                return
            }
            dataArr.forEach(element => {
                element.color = state.getClassColor(element.iaqclass)
            });

            let data = dataArr[dataArr.length - 1]
            state.last = {}
            let state1 = state.last
            state1.timeStamp = data.timeStamp
            state1.tempraw = data.tempraw
            state1.press = data.press
            state1.humiraw = data.humiraw
            state1.gaso = data.gaso
            state1.iaq = data.iaq
            state1.iaqacc = data.iaqacc
            state1.temp = data.temp
            state1.humy = data.humy
            state1.co2 = data.co2
            state1.voc = data.voc
            state1.sensorid = data.sensorid
            state1.place = data.place
        }
    },
    
}