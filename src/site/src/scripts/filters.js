app.filter('capitalize', function() {
    return function(input) {
        output = input[0].toUpperCase() + input.substr(1,input.length);
        return output;
    };
});
