Template.polls.helpers({
    polls: function() {
        return Polls.find({}, {sort: {createdAt: -1}});
    }
});

Template.polls.events({
    "click #newPollBtn": function() {
        Session.set("showNewPoll", !Session.get("showNewPoll"));
    }
});

Template.poll.created = function() {
    this.data.type = "poll";
};

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
