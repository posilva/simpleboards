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

### 
