htmx.defineExtension('json-enc', {
    onEvent: function (name, evt) {
        if (name === "htmx:configRequest") {
            evt.detail.headers['Content-Type'] = "application/json";
        }
    },

    /**
     * @param {FormData} parameters - [TODO:description]
     */
    encodeParameters : function(xhr, parameters, _) {
        xhr.overrideMimeType('text/json');
        let ob = {}

        parameters.keys().forEach(k => {
            const values = parameters.getAll(k)
            let key = k
            let array = false
            let int = false

            if (key.endsWith("[]")) {
                array = true
                key = key.slice(0, -2)
            }
            if (key.endsWith(".int")) {
                int = true
                key = key.slice(0, -4)
            }

            if (array) {
                ob[key] = []
                for (let i in values) {
                    ob[key].push(int ? ~~values[i] : values[i])
                }
            } else {
                ob[key] = int ? ~~values[0] : values[0]
            }
        })

        return JSON.stringify(ob);
    }
});
