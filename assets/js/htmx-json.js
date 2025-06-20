htmx.defineExtension('json-enc', (() => {
    const parseValue = (v, type) => {
        switch (type) {
        case "int":
            return parseInt(v)
        case "float":
            return parseFloat(v)
        case "bool":
            return !!v
        default:
            return v
        }
    }

    const buildObject  = (ob, params) => key => {
        const values = params.getAll(key)
        let type = 'str'
        let rest = []

        if (key.endsWith(".int")) {
            type = "int"
            key = key.slice(0, -4)
        } else if (key.endsWith(".float")) {
            type = "float"
            key = key.slice(0, -6)
        } else if (key.endsWith(".bool")) {
            type = "bool"
            key = key.slice(0, -5)
        }

        if (key.includes(".")) {
            [ key, ...rest ] = key.split(".")
        }

        if (values.length == 1 && !key.endsWith('[]')) {
            const v = parseValue(values[0], type)

            if (rest.length) {
                if (key in ob) {
                    ob[key][rest.join(".")] = v
                } else {
                    ob[key] = { [rest.join(".")]: v }
                }
            } else {
                ob[key] = v
            }

            return
        }

        if (key.endsWith('[]')) {
            key = key.slice(0, -2)
        }

        if (!(key in ob)) {
            ob[key] = []
        }

        for (let i in values) {
            const v = parseValue(values[i], type)

            if (!rest.length) {
                ob[key].push(v)
                continue
            }

            if (i < ob[key].length) {
                ob[key][i][rest.join(".")] = v
            } else {
                ob[key].push({ [rest.join(".")]: v })
            }
        }
    }

    return {
        onEvent: function (name, evt) {
            if (name === "htmx:configRequest") {
                evt.detail.headers['Content-Type'] = "application/json"
            }
        },
        encodeParameters : function(xhr, parameters, _) {
            xhr.overrideMimeType('text/json')
            let ob = {}

            console.log(JSON.stringify(parameters))
            parameters.keys().forEach(buildObject(ob, parameters))

            return JSON.stringify(ob)
        }
    }
})())
