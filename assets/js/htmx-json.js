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

        console.log(parameters.keys())
        parameters.keys().forEach(k => {
            console.log([k, parameters.get(k)])
            if (k.endsWith(".int")) {
                console.log(k.slice(0, -4), ~~parameters.get(k))
                ob[k.slice(0, -4)] = ~~parameters.get(k)
            } else {
                console.log("no", k, parameters.get(k))
                ob[k] = parameters.get(k)
            }
        })

        return JSON.stringify(ob);
    }
});
