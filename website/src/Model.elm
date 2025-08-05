module Model exposing (..)


init : () -> ( Model, Cmd msg )
init _ =
    ( { lgd =
            { rg = startingRawGame
            , legalMoves = []
            , selectedStartI = Nothing
            , boardFlipped = False
            , plr1blue = Human
            , plr2red = Ai Simple
            , futurePlr1blue = Human
            , futurePlr2red = Ai Simple
            }
      , tab = Local
      }
    , Cmd.none
    )


type alias Model =
    { lgd : LocalGameData
    , tab : Tab
    }


type Tab
    = Local
    | Online


type alias LocalGameData =
    { rg : RawGame -- see `RawGame` doc
    , legalMoves : List Move
    , selectedStartI : Maybe Int
    , boardFlipped : Bool

    -- the ones currently playing
    , plr1blue : Opponent
    , plr2red : Opponent

    -- the ones that will play once presses the play button
    , futurePlr1blue : Opponent
    , futurePlr2red : Opponent
    }


type Opponent
    = Human
    | Ai AiDifficulty


{-| this type is mirroring Game struct

    type Game struct {
        State GameState
        PlrTurn Player
        Board Board
        TurnNumber int
        TimeSinceExcitingMove int
    }

that way, Elm is the one keeping all the state,
and will pass it threw js to wasm(Go) when needs to make a move/ get the ai's move

    GameState:

    Playing GameState = "Playing"
    BlueWon GameState = "Blue Won"
    RedWon  GameState = "Red Won"
    Draw    GameState = "Draw"

    PlrTurn:

    0 = BluePlayer
    1 = RedPlayer

-}
type alias RawGame =
    { state : String
    , plrTurn : Int
    , board : List Int
    , turnNumber : Int
    , timeSinceExcitingMove : Int
    }


startingRawGame : RawGame
startingRawGame =
    { state = "Playing"
    , plrTurn = 0
    , board =
        [ [ 0, 2, 0, 2, 0, 2, 0, 2 ]
        , [ 2, 0, 2, 0, 2, 0, 2, 0 ]
        , [ 0, 2, 0, 2, 0, 2, 0, 2 ]
        , [ 1, 0, 1, 0, 1, 0, 1, 0 ]
        , [ 0, 1, 0, 1, 0, 1, 0, 1 ]
        , [ 4, 0, 4, 0, 4, 0, 4, 0 ]
        , [ 0, 4, 0, 4, 0, 4, 0, 4 ]
        , [ 4, 0, 4, 0, 4, 0, 4, 0 ]
        ]
            |> List.concat
    , timeSinceExcitingMove = 0
    , turnNumber = 1
    }


type alias Move =
    { startI : Int
    , endI : Int
    }


startingLegalMoves : List Move
startingLegalMoves =
    [ { endI = 24, startI = 17 }, { endI = 26, startI = 17 }, { endI = 26, startI = 19 }, { endI = 28, startI = 19 }, { endI = 28, startI = 21 }, { endI = 30, startI = 21 }, { endI = 30, startI = 23 } ]


type AiDifficulty
    = Easy
    | Medium
    | Hard
    | ExtraHard
    | Impossible
    | Simple


aiDifficulties : List AiDifficulty
aiDifficulties =
    [ Easy, Medium, Hard, ExtraHard, Impossible, Simple ]
