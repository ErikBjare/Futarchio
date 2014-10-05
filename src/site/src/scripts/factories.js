app.factory('msgStack', function() {
    var msgs = [];
    return function(obj) {
        if(msg === "") return msgs;
        msgs.push({type: "info", msg: msg});
        console.log(obj);
    };
});

app.factory('User', function($log, $resource) {
    var User = $resource('/api/0/:resource/:key/:val', {});
    return function(key, val) {
        var user = User.get(
           {"resource": "users",
            "key": key,
            "val": val},
            function(u) {
                $log.info(u.data[0]);
            }, function(u) {
                $log.error("Error");
            });
        return user;
    };
});
