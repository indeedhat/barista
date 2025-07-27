/**
 * Remap the json request body based on the following input name spec
 *
 * Cast to type
 * - Integer: suffix with .int
 * - Float: suffix with .float
 * - Boolean: suffix with .bool
 *
 * Cast as list:
 *   Suffix with []
 *   anything with the same prefix will be added to the same list
 *
 * Nested Lists:
 *   currently this is not supported
 *
 * Cast as typed list:
 *   Suffix with [] then the type ([].int)
 *
 * Cast to list of dicts:
 *   When a list has a .suffix the entry will be treat as a dictionary
 *   when encountering another entry with the same [] .suffix a new dictionary entry will be created
 *
 * Example:
 *   <input name="title" value="main title"/>
 *   <input name="is_decaf.bool" type="checkbox" value="1" checked />
 *   <input name="rating.float" value="7.5" />
 *   <input name="steps[].title" value="step one" />
 *   <input name="steps[].time.int" value="30" />
 *   <input name="steps[].title" value="step two" />
 *   <input name="steps[].time.int" value="45" />
 *
 *   Becomes:
 *   {
 *       "title": "main title",
 *       "is_decaf.bool": "1",
 *       "rating.float": "7.5",
 *       "steps[].title": "step one",
 *       "steps[].time.int": "30",
 *       "steps[].title": "step two"
 *       "steps[].time.int": "45"
 *   }
 *
 *   Produces:
 *   {
 *       "title": "main title",
 *       "is_decaf": true,
 *       "rating": 7.5,
 *       "steps": [
 *           {
 *               "title": "step one",
 *               "time": 30
 *           },
 *           {
 *               "title": "step two",
 *               "time": 30
 *           }
 *       ]
 *   }
 */
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
