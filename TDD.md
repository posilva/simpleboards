# Technical Design Document


### Leaderboard configuration 

```json
{
    "name": "fashion_game_weekly",
    "reset": "weekly",
    "function": "sum"
    "prizes_table:": [
        {
            "rank_from": 1,
            "rank_to": 3,
            "action": "gold"
        },
        {
            "rank_from": 4,
            "rank_to": 6,
            "action": "bronze"
        }
    ]
}

```

### On startup 

- Load in memory all leaderboards configs stored in the database  
- Monitor configs for periodic update 

### On Reset 

- Get the reset lock for a given epoch
- Process the prize rewards recording in the database who should receive the rewards 

### What's done 

- Set Score
- Get Current Epoch Leaderboards Scores

#### TODO 
- [ ] - Get Leaderboard Results based on the configuration

Leaderboards User Entry 
PK: USR#USER_ID
SK: LBRD#<name>::<epoch>

Queries:
- Return all the leaderboards epochs of an user
    - PK= USER#<user_id>, SK=beginwith(LBRD#<name>)

- Return all the leaderboards of an user
    - PK= USER#<user_id>, SK=beginwith(LBRD#)




Leaderboards Config
PK: LBRD#CONFIG
SK: LBRD#NAME#<name>







### Random notes

- when putting a score we can add some metadata that will influence how what data to store in the records and also the scoreboards will be storing the score
    - example:
        - a country based score board will be based on country code and will enforce that the country cannot change during an epoch
        - a league based score board will use a league data to put the score in the score board of a league
