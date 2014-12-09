Schemas = {};

// TODO: Remove?
Schemas.Score = new SimpleSchema({
    up: {
        type: Number,
        min: 0
    },
    down: {
        type: Number,
        min: 0
    }
});

Schemas.UserCreated = new SimpleSchema({
    createdBy: {
        type: String
    },
    createdAt: {
        type: Date
    }
});

Schemas.Post = new SimpleSchema([Schemas.UserCreated, {
    title: {
        type: String
    },
    description: {
        type: String
    }
}]);

Schemas.Poll = new SimpleSchema([Schemas.Post, {
    score: {
        type: Schemas.Score
    }
}]);

Schemas.Statement = new SimpleSchema([Schemas.Post, {
    score: {
        type: Schemas.Score
    }
}]);

Schemas.Prediction = new SimpleSchema([Schemas.UserCreated, {
    credence: {
        type: Number
    },
    statement: {
        type: String
    }
}]);

Schemas.Vote = new SimpleSchema({
    type: {
        type: String
    },
    value: {
        type: Number
    },
    post: {
        type: String
    }
});

function addCreationData(data) {
    data.createdBy = Meteor.userId();
    data.createdAt = new Date();
    return data;
}

Score = function() {
    return {"up": 0, "down": 0};
};

Polls = new Mongo.Collection("polls");
Polls.attachSchema(Schemas.Poll);
Poll = function(poll) {
    poll = addCreationData(poll);
    poll.score = new Score();
    return poll;
};

Statements = new Mongo.Collection("statements");
Statements.attachSchema(Schemas.Statement);
Statement = function(stmt) {
    stmt = addCreationData(stmt);
    stmt.score = new Score();
    return stmt;
};

Predictions = new Mongo.Collection("predictions");
Predictions.attachSchema(Schemas.Prediction);
Prediction = function(pred) {
    pred = addCreationData(pred);
    return pred;
};

Votes = new Mongo.Collection("votes");
Votes.attachSchema(Schemas.Vote);
Vote = function(vote) {
    vote.createdBy = Meteor.userId();
    vote.createdAt = new Date();
    return vote;
};
