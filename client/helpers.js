Template.registerHelper("session", function(str) {
    console.warn("Deprecated");
    return Session.get(str);
});

Template.registerHelper('eq', function(v1, v2, options) {
    if(v1 === v2){
        return true;
    } else {
        return false;
    }
});

Template.registerHelper('length', function(l) {
    return l.length;
});

Template.registerHelper("userId", function() {
    return Meteor.userId();
});

Template.registerHelper("usernameOf", function(userid) {
    var user = Meteor.users.findOne({"_id": userid});
    if(user === undefined) {
        console.error("User " + userid + " not found");
        return "[ERROR: USER NOT FOUND]";
    }
    if(user.username === undefined) {
        console.error("User " + userid + " has no username");
        return "?";
    }
    return user.username;
});

Template.registerHelper("capitalize", function(str) {
    if(str !== undefined && str.length > 0) {
        return str[0].toUpperCase() + str.substr(1);
    } else {
        console.warn("capitalize function got undefined or empty string as argument");
    }
});

Template.registerHelper("fromNow", function(date) {
    if(date === undefined) {
        console.error("date missing");
        return "";
    } else {
        return moment(date).fromNow();
    }
});

Template.registerHelper("tvar", function(varname) {
    return TemplateVar.get(varname);
});

Template.registerHelper("points", function(post) {
    var ratings = _.map(post.ratings, function(rating) {
        return rating.score;
    });
    var points = _.reduce(ratings, function(memo, rating) {
        return memo+rating;
    }, 0);
    return points.toString();
});

