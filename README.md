# socket-sweeper

Websocket minesweeper game. Watch other players live or stream your game out to the public. Watch previous leaderboard performances play by play.

## publisher ws route

Publisher can reconnect to original game if key matches one in publisher slice.

Publisher can create a new game if key does not match one in slice.

publisher can send actions to board resulting in game changes being send to subscribers.

### publisher
