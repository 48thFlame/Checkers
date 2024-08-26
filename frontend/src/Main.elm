module Main exposing (main)

import Browser
import Html
import Html.Attributes exposing (alt, class, src, style)
import Html.Events exposing (onClick, onInput)
import Translator exposing (..)


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


init : () -> ( Model, Cmd Msg )
init _ =
    ( { rg = defaultRawGame, difficulty = Simple }, Cmd.none )


type alias Model =
    { rg : RawGame
    , difficulty : AiDifficulty
    }


type Msg
    = UpdatedGameAppeared RawGame
    | MakeAction JsActions
    | ChangeDifficulty String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdatedGameAppeared rawGame ->
            ( { model | rg = rawGame }, Cmd.none )

        MakeAction action ->
            ( model, translator action )

        ChangeDifficulty stringedDiff ->
            ( { model | difficulty = stringToAiDifficulty stringedDiff }, Cmd.none )


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
            List.map intToSlot model.rg.board
    in
    Html.div [ class "main-area" ]
        [ boardToHtml displayBoard
        , Html.div [ class "control-area" ]
            [ Html.button [ onClick (MakeAction (GetAiMove model.difficulty model.rg)) ] [ Html.text "Play AI" ]
            , Html.button [ onClick (MakeAction GetNewGame) ] [ Html.text "New Game" ]
            , Html.select [ onInput ChangeDifficulty ]
                (List.map (aiDifficultyToHtmlOption model.difficulty) aiDifficulties)
            ]
        ]
