port module Main exposing (main)

import Browser
import Html
import Html.Attributes exposing (alt, class, src, style)
import Html.Events exposing (onClick)
import Json.Decode as JD
import Json.Encode as JE


main : Program JD.Value Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


init : JD.Value -> ( Model, Cmd Msg )
init flags =
    let
        rawGame =
            case decodeRawGame flags of
                Ok rg ->
                    rg

                Err _ ->
                    { state = "Playing"
                    , plrTurn = 0
                    , board = List.repeat 64 0
                    , timeSinceExcitingMove = 0
                    , turnNumber = 1
                    }
                        |> Debug.log "yeah failed.."
    in
    ( { rawGame = rawGame }, translator GetNewGame )


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

    type alias RawGame =
        { state : String
        , plrTurn : Int
        , board : List Int
        , turnNumber : Int
        , timeSinceExcitingMove : Int
        }

-}
type alias RawGame =
    { state : String
    , plrTurn : Int
    , board : List Int
    , turnNumber : Int
    , timeSinceExcitingMove : Int
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


type alias Model =
    { rawGame : RawGame
    }


type JsActions
    = GetNewGame
    | GetAiMove RawGame


translator : JsActions -> Cmd msg
translator action =
    case action of
        GetNewGame ->
            actionRequest ( "getNewGame", JE.null )

        GetAiMove rg ->
            let
                data =
                    JE.object
                        [ ( "game", encodeRawGame rg ) ]
            in
            actionRequest ( "getAiMove", data )


port actionRequest : ( String, JE.Value ) -> Cmd msg


port rawGameReceiver : (RawGame -> msg) -> Sub msg


type Msg
    = UpdatedGameAppeared RawGame
    | PlayAiMove
    | NewGame


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdatedGameAppeared rawGame ->
            ( { model | rawGame = rawGame }, Cmd.none )

        PlayAiMove ->
            ( model, translator (GetAiMove model.rawGame) )

        NewGame ->
            ( model, translator GetNewGame )


subscriptions : Model -> Sub Msg
subscriptions _ =
    rawGameReceiver UpdatedGameAppeared


intToSlot : Int -> Slot
intToSlot n =
    case n of
        0 ->
            NaS

        1 ->
            Empty

        2 ->
            BluePiece

        3 ->
            BlueKing

        4 ->
            RedPiece

        5 ->
            RedKing

        _ ->
            NaS


type Slot
    = NaS
    | Empty
    | BluePiece
    | BlueKing
    | RedPiece
    | RedKing


slotToHtml : Int -> Slot -> Html.Html msg
slotToHtml i slot =
    let
        slotISpan =
            Html.span [ class "slot-i" ]
                [ Html.text <| String.fromInt i ]

        slotImg pieceName =
            Html.img
                [ src ("assets/" ++ pieceName ++ ".webp")
                , alt pieceName
                , class "slot-img"
                ]
                []
    in
    case slot of
        NaS ->
            Html.div [ class "slot NaS-slot" ] []

        Empty ->
            Html.div [ class "slot empty-slot" ]
                [ slotISpan ]

        BluePiece ->
            Html.div [ class "slot piece-slot" ]
                [ slotISpan
                , slotImg "bluePiece"
                ]

        BlueKing ->
            Html.div [ class "slot piece-slot" ]
                [ slotISpan
                , slotImg "blueKing"
                ]

        RedPiece ->
            Html.div [ class "slot piece-slot" ]
                [ slotISpan
                , slotImg "redPiece"
                ]

        RedKing ->
            Html.div [ class "slot piece-slot" ]
                [ slotISpan
                , slotImg "redKing"
                ]


boardToHtml : List Slot -> Html.Html Msg
boardToHtml board =
    let
        stringedBoardSize =
            String.fromInt 8
    in
    Html.div
        [ class "checker-board"
        , style "grid-template-rows" ("repeat(" ++ stringedBoardSize ++ ", 1fr)")
        , style "grid-template-columns" ("repeat(" ++ stringedBoardSize ++ ", 1fr)")
        ]
        (List.indexedMap slotToHtml board)


view : Model -> Html.Html Msg
view model =
    let
        displayBoard =
            List.map intToSlot model.rawGame.board
    in
    Html.div [ class "main-area" ]
        [ boardToHtml displayBoard
        , Html.div [ class "control-area" ]
            [ Html.button [ onClick PlayAiMove ] [ Html.text "Play AI" ]
            , Html.button [ onClick NewGame ] [ Html.text "New Game" ]
            ]
        ]
