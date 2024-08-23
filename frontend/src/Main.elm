module Main exposing (main)

import Browser
import Html
import Html.Attributes exposing (alt, class, src, style)
import Json.Decode as JD


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
            case JD.decodeValue rawGameDecoder flags of
                Ok rg ->
                    rg

                Err _ ->
                    { state = 0
                    , plrTurn = 0
                    , board = List.repeat 64 0
                    , timeSinceExcitingMove = 0
                    , turnNumber = 1
                    }
    in
    ( { rawGame = rawGame }, Cmd.none )


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

-}
type alias RawGame =
    { state : Int
    , plrTurn : Int
    , board : List Int
    , turnNumber : Int
    , timeSinceExcitingMove : Int
    }


rawGameDecoder : JD.Decoder RawGame
rawGameDecoder =
    JD.map5 RawGame
        (JD.field "state" JD.int)
        (JD.field "plrTurn" JD.int)
        (JD.field "board" (JD.list JD.int))
        (JD.field "turnNumber" JD.int)
        (JD.field "timeSinceExcitingMove" JD.int)


type alias Model =
    { rawGame : RawGame
    }


type Msg
    = Msg1


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Msg1 ->
            ( model, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none


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
        [ boardToHtml displayBoard ]
