module View exposing (..)

import Html
import Html.Attributes exposing (alt, class, selected, src, style, value)
import Html.Events exposing (onClick, onInput)
import Model exposing (..)
import Translator exposing (..)
import Update exposing (..)


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


view : Model -> Html.Html Msg
view model =
    let
        _ =
            Debug.log "moves:" model.legalMoves

        displayBoard =
            List.map intToSlot model.rg.board
    in
    Html.div [ class "main-area" ]
        [ boardToHtml displayBoard
        , Html.div [ class "control-area" ]
            [ Html.button [ onClick (MakeAction (GetAiMove model.difficulty model.rg)) ] [ Html.text "Play AI" ]
            , Html.button [ onClick (UpdatedGameAppeared startingRawGame) ] [ Html.text "New Game" ]
            , Html.select [ onInput (\stringedDiff -> ChangeDifficulty (stringToAiDifficulty stringedDiff)) ]
                (List.map (aiDifficultyToHtmlOption model.difficulty) aiDifficulties)
            ]
        ]
