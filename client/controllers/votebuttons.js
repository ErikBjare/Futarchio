Template.votebuttons.events({
    "click #up": function(event, template) {
        vote = new Vote({
            type: "UpDown",
            value: 1,
            post: template.data._id
        });
        Votes.insert(vote);
    },
    "click #down": function(event, template) {
        vote = new Vote({
            type: "UpDown",
            value: -1,
            post: template.data._id
        });
        Votes.insert(vote);
    }
});

