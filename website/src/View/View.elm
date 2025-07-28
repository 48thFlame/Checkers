module View.View exposing (view)

import Html
import Html.Attributes exposing (attribute, class, for, href, id, selected, value)
import Html.Events exposing (onClick, onInput)
import Model exposing (AiDifficulty(..), Model, Opponent(..), aiDifficulties)
import Update exposing (Msg(..))
import View.Board exposing (viewBoard)


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
