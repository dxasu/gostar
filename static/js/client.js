(function () {
    var form = document.getElementById("form"),
        udata = getQueryStringArgs(),
        timer;


    //保存信息
    $('.toast-frame .next-btn').on('click', function () {
        var msg = onInputChange();
        if (!msg) {
            return;
        }
        form.submit();
    });

    clearInterval(timer);
    timer = setInterval(function () {
        onInputChange()
        clearInterval(timer);
    }, 1 * 1000);

})();

function onInputChange() {
    var email_reg = /^[\w\W]+@[\w\W]+$/,
    number_reg = /^[0-9]+$/,
    puts = document.getElementsByTagName("input");
    console.log(">>>>>>>>>>>>>>")
    for (var i = 0; i < puts.length; i++) {
        console.log(puts[i].type, puts[i].value , !number_reg.test(puts[i].value))
        if (puts[i].type == "text" && puts[i].value == "" || 
        puts[i].type == "email" && !email_reg.test(puts[i].value) ||
        puts[i].type == "number" && !number_reg.test(puts[i].value)) {
            $('.next-btn').css("background-color", "#E6E8EB");
            return false;
        }
    }
    console.log(puts)
    $('.next-btn').css("background-color", "#3CBAFF");
    return true
}

function getQueryStringArgs() {
    var qs = (location.search.length > 0 ? location.search.substring(1) : ''),
        args = {},
        items = qs.length ? qs.split('&') : [],
        item = null, name = null, value = null;
    for (var i = 0; i < items.length; i++) {
        item = items[i].split('=');
        name = decodeURIComponent(item[0]);
        value = decodeURIComponent(item[1]);

        if (name.length) {
            args[name] = value;
        }
    }
    return args;
}
