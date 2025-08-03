module View.View exposing (view)

import Html
import Html.Attributes exposing (class)
import Model exposing (Model)
import Update exposing (Msg(..))
import View.Board exposing (viewBoard)
import View.Control exposing (viewControl)


view : Model -> Html.Html Msg
view model =
    Html.div [ class "app" ]
        [ viewBoard model.lgd
        , viewControl model.lgd
        ]



-- Html.div [ class "app" ]
--     -- [ Html.h1 [] [ Html.text "Play Checkers!" ]
--     [ viewControl model
--     -- , Html.footer [ class "credits" ]
--     --     [ Html.p []
--     --         [ Html.text "Website and Bot made by "
--     --         , Html.a [ href "https://github.com/48thFlame" ]
--     --             [ Html.text "48thFlame" ]
--     --         , Html.text " "
--     --         , Html.a [ href "https://github.com/48thFlame/Checkers" ]
--     --             [ Html.text "Repo" ]
--     --         ]
--     --     ]
--     ]
