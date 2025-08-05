module View.Control exposing (viewControl)

import Html
import Html.Attributes exposing (attribute, class, for, id, selected, value)
import Html.Events exposing (onClick, onInput)
import Model exposing (AiDifficulty(..), Model, Opponent(..), Tab(..), aiDifficulties)
import Update exposing (Msg(..))


plrSelect : String -> Opponent -> String -> (String -> Msg) -> Html.Html Msg
plrSelect labelText currentOpp selectionHtmlId changeDiffMsg =
    let
        -- | maybe difficulty, if its human then should be `Nothing`
        md : Maybe AiDifficulty
        md =
            case currentOpp of
                Human ->
                    Nothing

                Ai diff ->
                    Just diff

        stringedDiff : AiDifficulty -> String
        stringedDiff diff =
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
    Html.div []
        [ Html.label [ for selectionHtmlId ] [ Html.text labelText ]
        , Html.br [] []
        , Html.select
            [ id selectionHtmlId, class "plr-select", onInput changeDiffMsg ]
            [ Html.option [ value "Human", selected (currentOpp == Human) ] [ Html.text "You" ]
            , Html.hr [] []
            , Html.optgroup
                [ attribute "label" "Ai Difficulties:" ]
                (List.map
                    (\diff ->
                        Html.option
                            [ value (stringedDiff diff)

                            -- either will match at some point, or means that "Human" is selected and matched there
                            , selected (Just diff == md)
                            ]
                            [ Html.text (stringedDiff diff) ]
                    )
                    aiDifficulties
                )
            ]
        ]


tab : Tab -> Html.Html Msg
tab t =
    let
        stringedTab =
            case t of
                Local ->
                    "Local"

                Online ->
                    "Online"
    in
    Html.button [ class "tab-switcher-button", onClick (ChangeTab t) ] [ Html.text stringedTab ]


viewControl : Model -> Html.Html Msg
viewControl model =
    let
        tabSwitcherArea =
            Html.div []
                [ tab Local
                , tab Online
                ]

        localButtonsArea =
            Html.div
                [ if model.tab == Local then
                    class "local-tab"

                  else
                    class "non-active-tab"
                ]
                [ plrSelect "Player 1:" model.lgd.futurePlr1blue "plr1-select-id" ChangePlr1
                , plrSelect "Player 2:" model.lgd.futurePlr2red "plr2-select-id" ChangePlr2
                , Html.button [ class "newGame-button", onClick NewGame ]
                    [ Html.text "Play!" ]
                ]
    in
    Html.div [ class "control" ]
        [ tabSwitcherArea
        , localButtonsArea
        , Html.button [ class "flipBoard-button", onClick FlipBoard ]
            [ Html.text "Flip Board" ]
        ]
