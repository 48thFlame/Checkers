module View exposing (..)

import Html
import Html.Attributes exposing (alt, attribute, class, for, href, id, selected, src, style, value)
import Html.Events exposing (onClick, onInput)
import Model exposing (..)
import Set
import Translator exposing (..)
import Update exposing (..)


{-| a slot in a checkers board
`MoveStarter` - is a slot that goes on top of a regular one such that it shows that can start moving from there.
`MoveEnder` is the same just for ending
-}
type BaseSlot
    = NaS
    | Empty
    | BluePiece
    | BlueKing
    | RedPiece
    | RedKing


type OuterSlot
    = MoveStarter BaseSlot
    | Selected BaseSlot
    | MoveEnder BaseSlot
    | Regular BaseSlot


baseSlotToHtml : Int -> BaseSlot -> Html.Html Msg
baseSlotToHtml i slot =
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
            Html.div [ class "empty" ] [ slotISpan ]

        BluePiece ->
            piece "bluePiece"

        BlueKing ->
            piece "blueKing"

        RedPiece ->
            piece "redPiece"

        RedKing ->
            piece "redKing"


outerSlotToHtml : Bool -> Int -> OuterSlot -> Html.Html Msg
outerSlotToHtml flipped i slot =
    let
        rowStr =
            (i // 8)
                + 1
                |> (if flipped then
                        \ri -> 9 - ri

                    else
                        identity
                   )
                |> String.fromInt

        colStr =
            modBy 8 i
                + 1
                |> (if flipped then
                        \ci -> 9 - ci

                    else
                        identity
                   )
                |> String.fromInt

        baseAttrs =
            [ style "grid-row-start" rowStr
            , style "grid-column-start" colStr
            , class "slot"
            ]
    in
    case slot of
        Regular s ->
            Html.div baseAttrs [ baseSlotToHtml i s ]

        MoveStarter s ->
            Html.div
                ([ class "startI-slot", onClick (StartSlotSelected i) ]
                    ++ baseAttrs
                )
                [ baseSlotToHtml i s ]

        Selected s ->
            Html.div
                ([ class "startI-slot selectedI-slot", onClick (StartSlotSelected i) ]
                    ++ baseAttrs
                )
                [ baseSlotToHtml i s ]

        MoveEnder s ->
            Html.div ([ class "endI-slot", onClick (EndSlotSelected i) ] ++ baseAttrs)
                [ Html.div [ class "endI-slot-circle" ] []
                , baseSlotToHtml i s
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

        baseSlotToOuterSlot i s =
            if Just i == model.selectedStartI then
                Selected s

            else if Set.member i startIs then
                MoveStarter s

            else if List.member i endIs then
                MoveEnder s

            else
                Regular s

        intToBaseSlot n =
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

        gameStateDiv =
            Html.div [ class "game-state", style "grid-area" "1/1/9/9" ]
                [ Html.p [ class "game-state-text" ]
                    [ Html.text
                        (case model.rg.state of
                            "Draw" ->
                                "Game ended in a draw.."

                            "Blue Won" ->
                                "Blue won the game!!"

                            "Red Won" ->
                                "Red won the game!!"

                            _ ->
                                ""
                        )
                    , Html.br [] []
                    , Html.br [] []
                    , Html.text "Pres \"Play\" button to play another game."
                    ]
                ]
    in
    Html.div
        [ class "checker-board"
        , style "grid-template-rows" "repeat(8, 1fr)"
        , style "grid-template-columns" "repeat(8, 1fr)"
        ]
        (model.rg.board
            |> List.map intToBaseSlot
            |> List.indexedMap baseSlotToOuterSlot
            |> List.indexedMap (outerSlotToHtml model.boardFlipped)
            |> (if model.rg.state /= "Playing" then
                    List.append [ gameStateDiv ]

                else
                    identity
               )
        )


aiDifficultyToHtmlOption : Maybe AiDifficulty -> AiDifficulty -> Html.Html msg
aiDifficultyToHtmlOption md diff =
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
    Html.option
        [ value stringedDiff, selected (Just diff == md) ]
        [ Html.text stringedDiff ]


plrSelectOptions : Opponent -> List (Html.Html Msg)
plrSelectOptions opp =
    let
        md =
            case opp of
                Human ->
                    Nothing

                Ai diff ->
                    Just diff
    in
    [ Html.option [ value "Human", selected (opp == Human) ] [ Html.text "You" ]
    , Html.hr [] []
    , Html.optgroup
        [ attribute "label" "Ai Difficulties" ]
        (List.map (aiDifficultyToHtmlOption md) aiDifficulties)
    ]


viewControlArea : Model -> Html.Html Msg
viewControlArea model =
    Html.div [ class "control-area" ]
        [ Html.label
            [ for "plr1-select", class "plr-select-label" ]
            [ Html.text "Player 1:" ]
        , Html.select
            [ id "plr1-select", class "ctrl-obj plr-select plr1-select", onInput ChangePlr1 ]
            (plrSelectOptions model.futurePlr1blue)
        , Html.label
            [ for "plr2-select", class "plr-select-label" ]
            [ Html.text "Player 2:" ]
        , Html.select
            [ id "plr2-select", class "ctrl-obj plr-select plr2-select", onInput ChangePlr2 ]
            (plrSelectOptions model.futurePlr2red)
        , Html.button
            [ class "ctrl-obj ctrl-button newGame-button", onClick NewGame ]
            [ Html.text "Play!" ]
        , Html.button
            [ class "ctrl-obj ctrl-button flipBoard-button", onClick FlipBoard ]
            [ Html.text "Flip Board" ]
        ]


view : Model -> Html.Html Msg
view model =
    Html.div [ class "app" ]
        [ Html.h1 [] [ Html.text "Play Checkers!" ]
        , viewBoard model
        , viewControlArea model
        , Html.div [ class "credits" ]
            [ Html.p []
                [ Html.text "Website and Bot made by "
                , Html.a [ href "https://github.com/48thFlame" ]
                    [ Html.text "48thFlame" ]
                , Html.text " "
                , Html.a [ href "https://github.com/48thFlame/Checkers" ]
                    [ Html.text "Repo" ]
                ]
            ]
        ]
