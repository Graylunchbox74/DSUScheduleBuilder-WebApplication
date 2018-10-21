modulejs.define('request_manager', function () {
    return {
        post: function (url, json_data) {
            return fetch(url, {
                method: "POST",
                cache: "no-cache",
                headers: {
                    "Content-Type": "application/json; charset=utf-8"
                },
                body: JSON.stringify(json_data)
            })
        },

        post_json: function (url, json_data) {
            return this.post(url, json_data).then(res => res.json())
        }
    }
});