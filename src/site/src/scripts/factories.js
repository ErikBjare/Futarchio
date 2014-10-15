app.factory('msgStack', function() {
    var msgs = [];
    return function(obj) {
        if(msg === "") return msgs;
        msgs.push({type: "info", msg: msg});
        console.log(obj);
    };
});

app.factory('User', function($log, $resource) {
    var User = $resource('/api/0/users/:key/:val', {});
    return function(key, val) {
        var user = User.query(
           {"key": key,
            "val": val},
            function(u) {
                $log.info("Successfully fetched user");
            }, function(u) {
                $log.error("Error");
            });
        return user;
    };
});

app.factory('gravatar', function($log, $resource) {
    var gravatar = {};
    gravatar.hash = function(email) {
        if(!_.contains(email, "@")) {
            console.error("Error: got empty email");
        }
        var hash = CryptoJS.MD5(email).toString();
        return hash;
    };

    return gravatar;
});

app.factory('user', function($q, $log, $http, $route, $cookieStore, $location, gravatar) {
    var user = {};
    user.is_logged_in = function() {
        val = $cookieStore.get("me") && $cookieStore.get("auth");
        return val;
    };

    user.getAuthkey = function() {
        authkey = $cookieStore.get("auth");
        return authkey;
    };

    user.login = function(username, password) {
        var deferred = $q.defer();

        $http.post('/api/0/auth', {username: username, password: password})
        .success(function(data, status, headers, config) {
            console.log(data);
            authkey = data.key;
            $cookieStore.put("auth", authkey);
            console.log("Cookie put");
            $http.get('/api/0/users/me', {"headers": {"Authorization": authkey}})
            .success(function(data) {
                console.log(data);
                $cookieStore.put("me", {"username": username, "email": data.email});
                deferred.resolve(data);
            }).error(function(data) {
                console.log("Error");
            });
        }).error(function(data, status, headers, config) {
            error = "Something went wrong while trying to make request";
            deferred.reject(error);
        });

        return deferred.promise;
    };

    user.logout = function() {
        $cookieStore.remove("auth");
        $cookieStore.remove("me");
        $location.path("/");

        // TODO: Use $window.location.reload() instead?
        $route.reload();
    };

    user.username = function() {
        return $cookieStore.get("me").username;
    };

    user.email = function() {
        return $cookieStore.get("me").email;
    };

    user.gravatar_hash = function() {
        return gravatar.hash(user.email());
    };

    return user;
});
