modulejs.define('request_manager', function () {
    return {
        post: function (url, json_data) {
            var fProm = fetch(url, {
                method: "POST",
                cache: "no-cache",
                headers: {
                    "Content-Type": "application/json; charset=utf-8"
                },
                body: JSON.stringify(json_data)
            })
            return new Promise((res, rej) => {
                fProm.then(data => {
                    // Check to see if we're being redirected to log in page
                    if (data.redirected == true && data.url.endsWith("/login")) {
                        window.location.pathname = "/logout"
                        rej("REDIRECT")
                        return
                    }
                    res(data)
                })
            })
        },

        post_json: function (url, json_data) {
            return this.post(url, json_data).then(res => res.json())
        }
    }
});