port module Translator exposing (..)

import Html
import Html.Attributes exposing (selected, value)
import Json.Decode as JD
import Json.Encode as JE


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


defaultRawGame : RawGame
defaultRawGame =
    { state = "Draw"
    , plrTurn = 0
    , board = List.repeat 64 0
    , timeSinceExcitingMove = 0
    , turnNumber = 1
    }


decodeRawGame : JD.Value -> Result JD.Error RawGame
decodeRawGame value =
    let
        decoder =
            JD.map5 RawGame
                (JD.field "state" JD.string)
                (JD.field "plrTurn" JD.int)
                (JD.field "board" (JD.list JD.int))
                (JD.field "turnNumber" JD.int)
                (JD.field "timeSinceExcitingMove" JD.int)
    in
    JD.decodeValue decoder value


encodeRawGame : RawGame -> JE.Value
encodeRawGame rawGame =
    JE.object
        [ ( "state", JE.string rawGame.state )
        , ( "plrTurn", JE.int rawGame.plrTurn )
        , ( "board", JE.list JE.int rawGame.board )
        , ( "turnNumber", JE.int rawGame.turnNumber )
        , ( "timeSinceExcitingMove", JE.int rawGame.timeSinceExcitingMove )
        ]


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


encodeAiDifficulty : AiDifficulty -> JE.Value
encodeAiDifficulty diff =
    JE.int
        (case diff of
            Easy ->
                1

            Medium ->
                2

            Hard ->
                3

            ExtraHard ->
                4

            Impossible ->
                5

            Simple ->
                6
        )


aiDifficultyToString : AiDifficulty -> String
aiDifficultyToString diff =
    case diff of
        Easy ->
            "Easy"

        Medium ->
            "Medium"

        Hard ->
            "Hard"

        ExtraHard ->
            "ExtraHard"

        Impossible ->
            "Impossible"

        Simple ->
            "Simple"


stringToAiDifficulty : String -> AiDifficulty
stringToAiDifficulty str =
    case str of
        "Easy" ->
            Easy

        "Medium" ->
            Medium

        "Hard" ->
            Hard

        "ExtraHard" ->
            ExtraHard

        "Impossible" ->
            Impossible

        "Simple" ->
            Simple

        _ ->
            Simple


aiDifficultyToHtmlOption : AiDifficulty -> AiDifficulty -> Html.Html msg
aiDifficultyToHtmlOption selectedDiff diff =
    let
        stringedDiff =
            aiDifficultyToString diff
    in
    Html.option [ value stringedDiff, selected (selectedDiff == diff) ] [ Html.text stringedDiff ]


type JsActions
    = GetNewGame
    | GetAiMove AiDifficulty RawGame


translator : JsActions -> Cmd msg
translator action =
    case action of
        GetNewGame ->
            actionRequest ( "getNewGame", JE.null )

        GetAiMove diff rg ->
            let
                data =
                    JE.object
                        [ ( "game", encodeRawGame rg )
                        , ( "difficulty", encodeAiDifficulty diff )
                        ]
            in
            actionRequest ( "getAiMove", data )


port actionRequest : ( String, JE.Value ) -> Cmd msg


port rawGameReceiver : (RawGame -> msg) -> Sub msg
