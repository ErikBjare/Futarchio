Template.polls.helpers({
    polls: function() {
        return Polls.find({}, {sort: {createdAt: -1}});
    }
});

Template.poll.helpers({
    error: ""
});

Template.statement.events({
    "click button#predict": function() {
        Predictions.insert({
            credence: event.target.credence.value,
            createdBy: Meteor.userId()
        });
        console.log("Inserted prediction");
    }
});

Template.statements.helpers({
    statements: function() {
        return Statements.find({}, {sort: {"createdAt": -1}});
    },
    showAdd: function() {
        return Session.get("addStmt-visible");
    }
});

Template.statements.events({
    "click button#addbtn": function(event) {
        Session.set("addStmt-visible", !Session.get("addStmt-visible"));
    }
});

Template.newpoll.events({
    "submit": function(event) {
        poll = new Poll({
            title: event.target.title.value,
            description: event.target.description.value
        });
        Polls.insert(poll);

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});

Template.newstatement.events({
    "submit": function(event) {
        stmt = new Statement({
            title: event.target.title.value,
            description: event.target.description.value
        });
        Statements.insert(stmt);

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});
