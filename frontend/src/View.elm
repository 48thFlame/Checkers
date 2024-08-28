module View exposing (..)

import Html
import Html.Attributes exposing (alt, class, selected, src, style, value)
import Html.Events exposing (onClick, onInput)
import Model exposing (..)
import Set
import Translator exposing (..)
import Update exposing (..)


{-| a slot in a checkers board
`MoveStarter` - is a slot that goes on top of a reglur one such that it shows that can start moving from there.
`MoveEnder` is the same just for ending
-}
type Slot
    = NaS
    | Empty
    | BluePiece
    | BlueKing
    | RedPiece
    | RedKing
    | MoveStarter
    | MoveEnder


slotToHtml : Int -> Slot -> Html.Html msg
slotToHtml i slot =
    let
        rowStr =
            (i // 8) + 1 |> String.fromInt

        colStr =
            modBy 8 i + 1 |> String.fromInt

        locStyles =
            [ style "grid-row-start" rowStr
            , style "grid-column-start" colStr
            ]

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

        pieceSlot name =
            Html.div (class "slot piece-slot" :: locStyles)
                [ slotISpan
                , slotImg name
                ]
    in
    case slot of
        NaS ->
            Html.div (class "slot NaS-slot" :: locStyles) []

        Empty ->
            Html.div (class "slot empty-slot" :: locStyles)
                [ slotISpan ]

        BluePiece ->
            pieceSlot "bluePiece"

        BlueKing ->
            pieceSlot "blueKing"

        RedPiece ->
            pieceSlot "redPiece"

        RedKing ->
            pieceSlot "redKing"

        MoveStarter ->
            Html.div (class "slot startI-slot" :: locStyles) []

        MoveEnder ->
            Html.div (class "slot endI-slot" :: locStyles) []


{-| convert list of slots to html, marking all location that are "startI" and "endI"
-}
viewBoard : Model -> Html.Html Msg
viewBoard model =
    let
        stringedBoardSize =
            "8"

        board =
            List.map
                (\n ->
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
                )
                model.rg.board

        startIs =
            List.map (\move -> move.startI) model.legalMoves
                |> Set.fromList
                -- remove duplicates
                |> Set.toList

        endIs =
            List.map (\move -> move.endI) model.legalMoves
                |> Set.fromList
                -- remove duplicates
                |> Set.toList

        startIsHtml =
            List.map (\i -> slotToHtml i MoveStarter) startIs

        endIsHtml =
            List.map (\i -> slotToHtml i MoveEnder) endIs
    in
    Html.div
        [ class "checker-board"
        , style "grid-template-rows" ("repeat(" ++ stringedBoardSize ++ ", 1fr)")
        , style "grid-template-columns" ("repeat(" ++ stringedBoardSize ++ ", 1fr)")
        ]
        (List.indexedMap slotToHtml board
            ++ startIsHtml
            ++ endIsHtml
        )


aiDifficultyToHtmlOption : AiDifficulty -> AiDifficulty -> Html.Html msg
aiDifficultyToHtmlOption selectedDiff diff =
    let
        stringedDiff =
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
    in
    Html.option [ value stringedDiff, selected (selectedDiff == diff) ] [ Html.text stringedDiff ]


view : Model -> Html.Html Msg
view model =
    Html.div [ class "app" ]
        [ viewBoard model
        , Html.div [ class "control-area" ]
            [ Html.button
                [ onClick (MakeAction (GetAiMove model.difficulty model.rg)) ]
                [ Html.text "Play AI" ]
            , Html.button
                [ onClick NewGame ]
                [ Html.text "New Game" ]
            , Html.select
                [ onInput ChangeDifficulty ]
                (List.map (aiDifficultyToHtmlOption model.difficulty) aiDifficulties)
            ]
        ]
