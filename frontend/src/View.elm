module View exposing (..)

import Html
import Html.Attributes exposing (alt, class, selected, src, style, value)
import Html.Events exposing (onClick, onInput)
import Model exposing (..)
import Set
import Translator exposing (..)
import Update exposing (..)


{-| a slot in a checkers board
`MoveStarter` - is a slot that goes on top of a regular one such that it shows that can start moving from there.
`MoveEnder` is the same just for ending
-}
type Slot
    = NaS
    | Empty
    | BluePiece
    | BlueKing
    | RedPiece
    | RedKing


type OuterSlot
    = MoveStarter Slot
    | Selected Slot
    | MoveEnder Slot
    | Regular Slot


innerSlotToHtml : Int -> Slot -> Html.Html Msg
innerSlotToHtml i slot =
    let
        slotISpan =
            Html.span [ class "slot-i" ]
                [ Html.text <| String.fromInt i ]

        slotImg pieceName =
            Html.img
                [ src ("assets/" ++ pieceName ++ ".webp")
                , alt pieceName
                , class "piece-img"
                ]
                []

        piece name =
            Html.div [ class "piece" ]
                [ slotISpan
                , slotImg name
                ]
    in
    case slot of
        NaS ->
            Html.div [ class "NaS" ] []

        Empty ->
            Html.div [ class "empty" ]
                [ slotISpan ]

        BluePiece ->
            piece "bluePiece"

        BlueKing ->
            piece "blueKing"

        RedPiece ->
            piece "redPiece"

        RedKing ->
            piece "redKing"


outerSlotToHtml : Int -> OuterSlot -> Html.Html Msg
outerSlotToHtml i slot =
    let
        rowStr =
            (i // 8) + 1 |> String.fromInt

        colStr =
            modBy 8 i + 1 |> String.fromInt

        baseAttrs =
            [ style "grid-row-start" rowStr
            , style "grid-column-start" colStr
            , class "slot"
            ]
    in
    case slot of
        Regular s ->
            Html.div baseAttrs [ innerSlotToHtml i s ]

        MoveStarter s ->
            Html.div
                ([ class "startI-slot", onClick (SlotSelected i) ]
                    ++ baseAttrs
                )
                [ innerSlotToHtml i s ]

        Selected s ->
            Html.div
                ([ class "startI-slot selectedI-slot", onClick (SlotSelected i) ]
                    ++ baseAttrs
                )
                [ innerSlotToHtml i s ]

        MoveEnder s ->
            Html.div (class "endI-slot" :: baseAttrs)
                [ Html.div [ class "endI-slot-circle" ] []
                , innerSlotToHtml i s
                ]


{-| convert list of slots to html, marking all location that are "startI" and "endI"
-}
viewBoard : Model -> Html.Html Msg
viewBoard model =
    let
        startIs =
            List.map (\move -> move.startI) model.legalMoves
                -- remove duplicates (necessary?)
                |> Set.fromList

        endIs =
            case model.selectedStartI of
                Nothing ->
                    []

                Just si ->
                    List.filter (\move -> move.startI == si) model.legalMoves
                        |> List.map (\move -> move.endI)

        innerSlotToOuter i s =
            if Just i == model.selectedStartI then
                Selected s

            else if Set.member i startIs then
                MoveStarter s

            else if List.member i endIs then
                MoveEnder s

            else
                Regular s

        slots =
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
                |> List.indexedMap innerSlotToOuter
    in
    Html.div
        [ class "checker-board"
        , style "grid-template-rows" "repeat(8, 1fr)"
        , style "grid-template-columns" "repeat(8, 1fr)"
        ]
        (List.indexedMap outerSlotToHtml slots)


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
                [ onClick (MakeAction (GetAiMove model.rg model.difficulty)) ]
                [ Html.text "Play AI" ]
            , Html.button
                [ onClick NewGame ]
                [ Html.text "New Game" ]
            , Html.select
                [ onInput ChangeDifficulty ]
                (List.map (aiDifficultyToHtmlOption model.difficulty) aiDifficulties)
            ]
        ]
