Template.registerHelper("session", function(str) {
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
        console.warning("capitalize function got undefined or empty string as argument");
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

Template.registerHelper("points", function(id) {
    var upVotes = Votes.find({"post": id, value: 1, type: "UpDown"}).count();
    var downVotes = Votes.find({"post": id, value: -1, type: "UpDown"}).count();
    var points = upVotes - downVotes;
    return points.toString();
});

