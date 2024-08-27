port module Translator exposing (..)

{-| Named translator because talks between go-js-elm
-}

import Json.Decode as JD
import Json.Encode as JE
import Model exposing (..)


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


type JsActions
    = GetAiMove AiDifficulty RawGame
    | GetLegalMoves RawGame


translator : JsActions -> Cmd msg
translator action =
    case action of
        GetAiMove diff rg ->
            let
                data =
                    JE.object
                        [ ( "game", encodeRawGame rg )
                        , ( "difficulty", encodeAiDifficulty diff )
                        ]
            in
            actionRequest ( "getAiMove", data )

        GetLegalMoves rg ->
            actionRequest ( "getLegalMoves", encodeRawGame rg )


port actionRequest : ( String, JE.Value ) -> Cmd msg


port rawGameReceiver : (RawGame -> msg) -> Sub msg


port legalMovesReceiver : (List Move -> msg) -> Sub msg
